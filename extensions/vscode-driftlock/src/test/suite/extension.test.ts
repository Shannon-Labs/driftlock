import * as assert from 'assert'
import { suite, test } from 'mocha'
import * as vscode from 'vscode'

suite('Extension Activation', () => {
  test('driftlock commands are registered', async () => {
    try {
      await vscode.commands.executeCommand('driftlock.startLiveLinting')
    } catch {
      // ignored for activation smoke test
    }
    const commands = await vscode.commands.getCommands(true)
    assert.ok(commands.includes('driftlock.startLiveLinting'))
    assert.ok(commands.includes('driftlock.stopLiveLinting'))
  })
})
