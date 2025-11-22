import * as vscode from 'vscode'
import { AnalyzerClient, AnalyzerSettings } from './analyzerClient'

export class TerminalTap implements vscode.Disposable {
  private terminal: vscode.Terminal | undefined
  private listener: vscode.Disposable | undefined
  private client: AnalyzerClient | undefined
  private buffer = ''

  constructor(
    private readonly getSettings: () => AnalyzerSettings,
    private readonly output: vscode.OutputChannel,
    private readonly updateStatus: (text: string) => void,
  ) {}

  async start() {
    if (this.listener) {
      vscode.window.showInformationMessage('Driftlock is already monitoring a terminal session.')
      return
    }
    const terminals = vscode.window.terminals
    if (terminals.length === 0) {
      vscode.window.showWarningMessage('No VS Code terminals are currently open.')
      return
    }
    const pick = await vscode.window.showQuickPick(terminals.map(t => ({ label: t.name, terminal: t })), {
      placeHolder: 'Select the terminal to monitor with Driftlock',
    })
    if (!pick) {
      return
    }

    const settings = this.getSettings()
    this.client = new AnalyzerClient(settings, this.output)
    await this.client.start()

    this.client.onResult(res => {
      if (res.is_anomaly) {
        this.output.appendLine(`[terminal:${pick.terminal.name}] score=${res.score.toFixed(2)} :: ${truncate(res.line)}`)
        this.updateStatus(`Terminal ALERT ${res.score.toFixed(2)}`)
      } else {
        this.updateStatus(`Terminal OK ${res.score.toFixed(2)}`)
      }
    })

    this.client.onError(err => this.output.appendLine(`Analyzer error: ${err.message}`))

    const api = (vscode.window as any)
    const onDidWrite = api.onDidWriteTerminalData as (callback: (event: any) => any) => vscode.Disposable | undefined
    if (!onDidWrite) {
      vscode.window.showErrorMessage('This version of VS Code does not support terminal data streaming APIs.')
      await this.stop()
      return
    }

    this.terminal = pick.terminal
    this.listener = onDidWrite((event: { terminal: vscode.Terminal; data: string }) => {
      if (!this.terminal || event.terminal !== this.terminal || !this.client) {
        return
      }
      this.buffer += event.data
      let newlineIndex = this.buffer.indexOf('\n')
      while (newlineIndex >= 0) {
        const line = this.buffer.slice(0, newlineIndex)
        this.buffer = this.buffer.slice(newlineIndex + 1)
        if (line.trim().length > 0) {
          this.client.write(line)
        }
        newlineIndex = this.buffer.indexOf('\n')
      }
    })

    vscode.window.showInformationMessage(`Driftlock monitoring terminal "${pick.terminal.name}"`)
  }

  async stop() {
    this.listener?.dispose()
    this.listener = undefined
    this.terminal = undefined
    this.buffer = ''
    if (this.client) {
      this.client.dispose()
      this.client = undefined
    }
    this.updateStatus('Terminal idle')
  }

  dispose() {
    void this.stop()
  }
}

function truncate(value: string, max = 160) {
  if (value.length <= max) {
    return value
  }
  return `${value.slice(0, max)}…`
}
