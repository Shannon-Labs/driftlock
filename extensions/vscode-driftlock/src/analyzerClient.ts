import * as vscode from 'vscode'
import { ChildProcessWithoutNullStreams, spawn } from 'child_process'

export interface AnalyzerSettings {
  cliPath: string
  baselineLines: number
  threshold: number
  format: 'raw' | 'ndjson'
  algo: 'zstd' | 'gzip'
  minLineLength: number
}

export interface AnalyzerResult {
  sequence: number
  line: string
  entropy: number
  baseline_entropy: number
  compression_ratio: number
  baseline_compression_ratio: number
  entropy_delta: number
  compression_delta: number
  score: number
  is_anomaly: boolean
  ready: boolean
  reason: string
}

export class AnalyzerClient implements vscode.Disposable {
  private child: ChildProcessWithoutNullStreams | undefined
  private stdoutBuffer = ''
  private readonly resultEmitter = new vscode.EventEmitter<AnalyzerResult>()
  private readonly errorEmitter = new vscode.EventEmitter<Error>()
  readonly config: AnalyzerSettings

  readonly onResult = this.resultEmitter.event
  readonly onError = this.errorEmitter.event

  constructor(settings: AnalyzerSettings, private readonly output: vscode.OutputChannel) {
    this.config = settings
  }

  async start(): Promise<void> {
    if (this.child) {
      return
    }

    const args = [
      'scan',
      '--format', this.config.format,
      '--baseline-lines', String(this.config.baselineLines),
      '--threshold', String(this.config.threshold),
      '--algo', this.config.algo,
      '--min-line-length', String(this.config.minLineLength),
      '--output', 'ndjson',
      '--show-all',
      '--stdin',
    ]

    await new Promise<void>((resolve, reject) => {
      this.child = spawn(this.config.cliPath, args, { stdio: 'pipe' })
      let resolved = false
      const finishResolve = () => {
        if (!resolved) {
          resolved = true
          resolve()
        }
      }

      this.child.once('error', err => {
        if (!resolved) {
          resolved = true
          reject(err)
        } else {
          this.errorEmitter.fire(err)
        }
      })

      this.child.stdout.on('data', chunk => this.consumeStdout(chunk.toString()))
      this.child.stderr.on('data', chunk => this.output.appendLine(`[driftlock-cli] ${chunk}`))

      this.child.on('exit', code => {
        this.child = undefined
        this.output.appendLine(`driftlock scan exited with code ${code ?? 'null'}`)
      })

      // Resolve on next tick after listeners attached
      setImmediate(finishResolve)
    })
  }

  write(line: string) {
    if (!this.child) {
      throw new Error('Analyzer process not started')
    }
    this.child.stdin.write(line + '\n')
  }

  dispose() {
    if (this.child) {
      this.child.stdin.end()
      this.child.kill()
      this.child = undefined
    }
    this.resultEmitter.dispose()
    this.errorEmitter.dispose()
  }

  private consumeStdout(chunk: string) {
    this.stdoutBuffer += chunk
    let newlineIndex = this.stdoutBuffer.indexOf('\n')
    while (newlineIndex >= 0) {
      const line = this.stdoutBuffer.slice(0, newlineIndex).trim()
      this.stdoutBuffer = this.stdoutBuffer.slice(newlineIndex + 1)
      if (line) {
        try {
          const payload = JSON.parse(line) as AnalyzerResult
          this.resultEmitter.fire(payload)
        } catch (err) {
          this.errorEmitter.fire(err as Error)
        }
      }
      newlineIndex = this.stdoutBuffer.indexOf('\n')
    }
  }
}
