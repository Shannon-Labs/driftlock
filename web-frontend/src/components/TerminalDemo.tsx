import { useState, useEffect } from "react";
import { Terminal, Play, Copy, Check } from "lucide-react";
import { Button } from "./ui/button";
import { Card } from "./ui/card";

const commands = [
  {
    command: "curl -X POST https://api.driftlock.net/api/v1/events/ingest",
    input: `{
  "type": "log",
  "source": "user-service",
  "data": { "level": "error", "message": "DB timeout" }
}`,
    output: `{
  "event_id": "evt_7f8a9b2c",
  "anomaly_detected": true,
  "compression_metrics": {
    "baseline": 3.2x,
    "current": 1.4x
  }
}`,
    delay: 1000,
  },
  {
    command: "curl https://api.driftlock.net/api/v1/anomalies/evt_7f8a9b2c",
    input: ``,
    output: `{
  "severity": "high",
  "explanation": "Compression dropped 56% — message field shows unusual pattern",
  "affected_fields": ["message"],
  "confidence": 0.85
}`,
    delay: 800,
  },
];

export const TerminalDemo = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [isTyping, setIsTyping] = useState(false);
  const [displayedCommand, setDisplayedCommand] = useState("");
  const [displayedInput, setDisplayedInput] = useState("");
  const [displayedOutput, setDisplayedOutput] = useState("");
  const [showOutput, setShowOutput] = useState(false);
  const [copied, setCopied] = useState(false);
  const [isPlaying, setIsPlaying] = useState(true);

  useEffect(() => {
    if (!isPlaying) return;

    const currentCommand = commands[currentStep];
    let charIndex = 0;

    setIsTyping(true);
    setShowOutput(false);
    setDisplayedCommand("");
    setDisplayedInput("");
    setDisplayedOutput("");

    // Type command
    const commandInterval = setInterval(() => {
      if (charIndex < currentCommand.command.length) {
        setDisplayedCommand(currentCommand.command.slice(0, charIndex + 1));
        charIndex++;
      } else {
        clearInterval(commandInterval);
        
        // Type input if exists
        if (currentCommand.input) {
          setTimeout(() => {
            let inputCharIndex = 0;
            const inputInterval = setInterval(() => {
              if (inputCharIndex < currentCommand.input.length) {
                setDisplayedInput(currentCommand.input.slice(0, inputCharIndex + 1));
                inputCharIndex++;
              } else {
                clearInterval(inputInterval);
                showResponse();
              }
            }, 5);
          }, 300);
        } else {
          showResponse();
        }
      }
    }, 30);

    const showResponse = () => {
      setTimeout(() => {
        setShowOutput(true);
        let outputCharIndex = 0;
        const outputInterval = setInterval(() => {
          if (outputCharIndex < currentCommand.output.length) {
            setDisplayedOutput(currentCommand.output.slice(0, outputCharIndex + 1));
            outputCharIndex++;
          } else {
            clearInterval(outputInterval);
            setIsTyping(false);
            
            // Move to next command after delay
            setTimeout(() => {
              setCurrentStep((prev) => (prev + 1) % commands.length);
            }, 3000);
          }
        }, 10);
      }, currentCommand.delay);
    };

    return () => {
      setIsTyping(false);
    };
  }, [currentStep, isPlaying]);

  const handleCopy = () => {
    const fullCommand = commands[currentStep].command + 
      (commands[currentStep].input ? "\n" + commands[currentStep].input : "");
    navigator.clipboard.writeText(fullCommand);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const togglePlay = () => {
    setIsPlaying(!isPlaying);
  };

  return (
    <Card className="relative glass-card rounded-2xl overflow-hidden shadow-2xl border-primary/20">
      {/* Terminal Header */}
      <div className="bg-muted/50 px-4 py-3 border-b border-border/40 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Terminal className="w-5 h-5 text-primary" />
          <span className="text-sm font-mono font-semibold">Driftlock API Demo</span>
        </div>
        <div className="flex items-center gap-2">
          <Button
            size="sm"
            variant="ghost"
            className="h-8 w-8 p-0"
            onClick={togglePlay}
          >
            <Play className={`w-4 h-4 ${isPlaying ? 'text-primary' : 'text-muted-foreground'}`} />
          </Button>
          <Button
            size="sm"
            variant="ghost"
            className="h-8 w-8 p-0"
            onClick={handleCopy}
          >
            {copied ? (
              <Check className="w-4 h-4 text-green-500" />
            ) : (
              <Copy className="w-4 h-4" />
            )}
          </Button>
        </div>
      </div>

      {/* Terminal Content */}
      <div className="bg-background/95 p-6 font-mono text-sm min-h-[500px] max-h-[600px] overflow-auto">
        {/* Command */}
        <div className="mb-4">
          <span className="text-primary">$</span>{" "}
          <span className="text-foreground">{displayedCommand}</span>
          {isTyping && !displayedInput && !showOutput && (
            <span className="inline-block w-2 h-4 bg-primary ml-1 animate-pulse"></span>
          )}
        </div>

        {/* Input JSON */}
        {displayedInput && (
          <div className="mb-4 ml-4">
            <pre className="text-muted-foreground whitespace-pre-wrap">
              {displayedInput}
            </pre>
            {isTyping && !showOutput && (
              <span className="inline-block w-2 h-4 bg-primary ml-1 animate-pulse"></span>
            )}
          </div>
        )}

        {/* Output */}
        {showOutput && (
          <div className="mt-4 mb-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
              <span className="text-xs text-green-500 font-semibold">RESPONSE</span>
            </div>
            <pre className="text-green-400 whitespace-pre-wrap ml-4">
              {displayedOutput}
            </pre>
            {isTyping && (
              <span className="inline-block w-2 h-4 bg-green-500 ml-4 animate-pulse"></span>
            )}
          </div>
        )}

        {/* Step Indicator */}
        <div className="flex gap-2 mt-8">
          {commands.map((_, idx) => (
            <div
              key={idx}
              className={`h-1 flex-1 rounded-full transition-colors ${
                idx === currentStep ? "bg-primary" : "bg-muted"
              }`}
            />
          ))}
        </div>
      </div>

      {/* Glow Effect */}
      <div className="absolute inset-0 bg-gradient-primary opacity-5 blur-3xl pointer-events-none"></div>
    </Card>
  );
};
