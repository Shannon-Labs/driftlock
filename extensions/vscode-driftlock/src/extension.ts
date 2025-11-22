import * as vscode from 'vscode'
import { AnalyzerClient, AnalyzerResult, AnalyzerSettings } from './analyzerClient'
import { TerminalTap } from './terminalTap'

let session: LiveLintSession | undefined
let terminalTap: TerminalTap | undefined

export function activate(context: vscode.ExtensionContext) {
  const diagnostics = vscode.languages.createDiagnosticCollection('driftlock')
  const decoration = vscode.window.createTextEditorDecorationType({
    backgroundColor: 'rgba(255, 0, 0, 0.25)'
  })
  const statusItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 10)
  statusItem.text = 'Δ Entropy: idle'
  statusItem.show()
  const output = vscode.window.createOutputChannel('Driftlock Live Radar')

  const getSettings = (): AnalyzerSettings => {
    const config = vscode.workspace.getConfiguration('driftlock')
    return {
      cliPath: config.get<string>('cliPath', 'driftlock'),
      baselineLines: config.get<number>('baselineLines', 400),
      threshold: config.get<number>('threshold', 0.35),
      format: config.get<'raw' | 'ndjson'>('format', 'raw'),
      algo: config.get<'zstd' | 'gzip'>('algo', 'zstd'),
      minLineLength: config.get<number>('minLineLength', 12)
    }
  }

  const startLiveLinting = async () => {
    if (session) {
      vscode.window.showInformationMessage('Driftlock live linting is already running.')
      return
    }
    const editor = vscode.window.activeTextEditor
    if (!editor) {
      vscode.window.showWarningMessage('Open a log file before starting Driftlock live linting.')
      return
    }
    const doc = editor.document
    const client = new AnalyzerClient(getSettings(), output)
    session = new LiveLintSession(doc, editor, client, diagnostics, decoration, statusItem, output, () => {
      session = undefined
    })
    try {
      await session.start()
    } catch (err) {
      session.dispose()
      session = undefined
      vscode.window.showErrorMessage(`Failed to start Driftlock: ${(err as Error).message}`)
    }
  }

  const stopLiveLinting = () => {
    if (!session) {
      vscode.window.showInformationMessage('Driftlock live linting is not running.')
      return
    }
    session.dispose()
    session = undefined
    statusItem.text = 'Δ Entropy: idle'
  }

  terminalTap = new TerminalTap(getSettings, output, text => {
    statusItem.text = `Δ Entropy: ${text}`
  })

  context.subscriptions.push(
    diagnostics,
    decoration,
    statusItem,
    output,
    vscode.commands.registerCommand('driftlock.startLiveLinting', startLiveLinting),
    vscode.commands.registerCommand('driftlock.stopLiveLinting', stopLiveLinting),
    vscode.commands.registerCommand('driftlock.startTerminalRadar', () => terminalTap?.start()),
    vscode.commands.registerCommand('driftlock.stopTerminalRadar', () => terminalTap?.stop())
  )
}

export function deactivate() {
  session?.dispose()
  terminalTap?.dispose()
}

class LiveLintSession implements vscode.Disposable {
  private readonly lineQueue: number[] = []
  private diagMap = new Map<number, vscode.Diagnostic>()
  private listener: vscode.Disposable | undefined
  private resultSub: vscode.Disposable | undefined
  private errorSub: vscode.Disposable | undefined
  private disposed = false
  private trackedLineCount: number
  private readonly settings: AnalyzerSettings

  constructor(
    private readonly document: vscode.TextDocument,
    private readonly editor: vscode.TextEditor,
    private readonly client: AnalyzerClient,
    private readonly diagnostics: vscode.DiagnosticCollection,
    private readonly decoration: vscode.TextEditorDecorationType,
    private readonly status: vscode.StatusBarItem,
    private readonly output: vscode.OutputChannel,
    private readonly onStop: () => void,
  ) {
    this.trackedLineCount = document.lineCount
    this.settings = client.config
  }

  async start() {
    await this.client.start()
    this.resultSub = this.client.onResult(res => this.handleResult(res))
    this.errorSub = this.client.onError(err => this.output.appendLine(`Analyzer error: ${err.message}`))
    this.primeBaseline()
    this.listener = vscode.workspace.onDidChangeTextDocument(event => this.onDocumentChange(event))
    this.status.text = 'Δ Entropy: monitoring'
  }

  dispose() {
    if (this.disposed) {
      return
    }
    this.disposed = true
    this.listener?.dispose()
    this.resultSub?.dispose()
    this.errorSub?.dispose()
    this.client.dispose()
    this.diagnostics.delete(this.document.uri)
    this.editor.setDecorations(this.decoration, [])
    this.onStop()
  }

  private primeBaseline() {
    const baselineStart = Math.max(0, this.document.lineCount -  this.settings.baselineLines)
    for (let line = baselineStart; line < this.document.lineCount; line++) {
      this.enqueueLine(line, this.document.lineAt(line).text)
    }
    this.trackedLineCount = this.document.lineCount
  }

  private enqueueLine(lineNumber: number, text: string) {
    this.lineQueue.push(lineNumber)
    this.client.write(text)
  }

  private onDocumentChange(event: vscode.TextDocumentChangeEvent) {
    if (event.document.uri.toString() !== this.document.uri.toString()) {
      return
    }
    const hasMidFileMutation = event.contentChanges.some(change => change.range.start.line < this.trackedLineCount - 1)
    if (hasMidFileMutation || event.document.lineCount < this.trackedLineCount) {
      vscode.window.showWarningMessage('Detected edits inside the monitored region. Restart Driftlock live linting to rebuild the baseline.')
      this.dispose()
      return
    }
    for (let line = this.trackedLineCount; line < event.document.lineCount; line++) {
      this.enqueueLine(line, event.document.lineAt(line).text)
    }
    this.trackedLineCount = event.document.lineCount
  }

  private handleResult(res: AnalyzerResult) {
    const lineNumber = this.lineQueue.shift()
    if (lineNumber === undefined) {
      return
    }
    if (!res.ready) {
      this.status.text = 'Δ Entropy: warming'
      return
    }
    this.status.text = `Δ Entropy: ${res.score.toFixed(2)} ${res.is_anomaly ? 'ALERT' : 'OK'}`
    if (!res.is_anomaly) {
      return
    }
    const line = this.document.lineAt(lineNumber)
    const range = new vscode.Range(lineNumber, 0, lineNumber, line.text.length)
    const diagnostic = new vscode.Diagnostic(range, res.reason || 'Entropy variance detected', vscode.DiagnosticSeverity.Warning)
    this.diagMap.set(lineNumber, diagnostic)
    this.flushDiagnostics()
  }

  private flushDiagnostics() {
    const list = Array.from(this.diagMap.values())
    this.diagnostics.set(this.document.uri, list)
    const ranges = list.map(diag => diag.range)
    this.editor.setDecorations(this.decoration, ranges)
  }
}
