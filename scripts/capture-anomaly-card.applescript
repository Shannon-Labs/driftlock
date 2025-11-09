-- Open demo HTML in Safari, scroll to first anomaly, capture window screenshot
set targetFile to POSIX file "/Volumes/VIXinSSD/driftlock/demo-output.html"

tell application "Safari"
  activate
  set docURL to "file:///Volumes/VIXinSSD/driftlock/demo-output.html"
  if (count of documents) is 0 then
    make new document with properties {URL:docURL}
  else
    set URL of front document to docURL
  end if
end tell

delay 1.5

tell application "Safari"
  try
    do JavaScript "document.querySelector('.anomaly-card')?.scrollIntoView({behavior:'instant', block:'center'});" in document 1
  end try
  tell window 1
    set bounds to {20, 80, 1200, 900}
  end tell
end tell

delay 0.8

set outPath to "/Volumes/VIXinSSD/driftlock/screenshots/demo-anomaly-card.png"

-- Try window capture via AX window id
set winId to ""
tell application "System Events"
  tell application process "Safari"
    set frontmost to true
    try
      set winId to value of attribute "AXWindowNumber" of window 1 as text
    end try
  end tell
end tell

try
  if winId is not "" then
    do shell script "/usr/sbin/screencapture -x -l " & winId & space & quoted form of outPath
  else
    do shell script "/usr/sbin/screencapture -x " & quoted form of outPath
  end if
end try

return outPath

