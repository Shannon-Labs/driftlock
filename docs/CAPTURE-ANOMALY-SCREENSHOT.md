# Step-by-Step Guide: Capture Anomaly Screenshot

This guide walks you through running the Driftlock demo and capturing a screenshot of the first anomaly card. Follow these steps in order.

---

## Part 1: Run the Demo

### Step 1: Navigate to the Project Directory
Open your terminal and type:
```bash
cd /Volumes/VIXinSSD/driftlock
```
Press Enter.

**What this does:** Changes your current directory to the Driftlock project folder.

---

### Step 2: Clean Up Old Files
Type this command:
```bash
rm -f demo-output.html driftlock-demo
```
Press Enter.

**What this does:** Removes any old demo files so you start fresh.

---

### Step 3: Build the Demo Program
Type this command:
```bash
go build -o driftlock-demo cmd/demo/main.go
```
Press Enter.

**What to expect:** The command will run and create a file called `driftlock-demo`. You should see a prompt again when it's done (no error messages).

**What this does:** Compiles the Go code into an executable program.

---

### Step 4: Run the Demo
Type this command:
```bash
./driftlock-demo test-data/financial-demo.json
```
Press Enter.

**What to expect:** You should see output like "Found 30 anomalies..." and a new file `demo-output.html` will be created.

**What this does:** Runs the demo program and generates an HTML report with anomaly results.

---

## Part 2: Open and Find the First Anomaly

### Step 5: Open the HTML File in Your Browser

**Option A - Safari (macOS default):**
```bash
open demo-output.html
```

**Option B - Chrome:**
```bash
open -a "Google Chrome" demo-output.html
```

**What to expect:** Your browser will open and show the demo results page.

---

### Step 6: Find the First Anomaly Card

1. In your browser, press **Command-F** (or **Ctrl-F** on Windows/Linux) to open the search box.
2. Type: `ANOMALY DETECTED`
3. Press Enter.

**What to expect:** The browser will highlight the first occurrence of "ANOMALY DETECTED" and scroll to it. This is your first anomaly card (it will have a red badge).

---

## Part 3: Capture the Screenshot

Choose ONE of these three methods. **Option B (Chrome DevTools) is the most precise.**

---

### Option A: macOS Selection Screenshot (Easiest)

1. Make sure the anomaly card is visible on your screen.
2. Press **Shift-Command-4** (hold all three keys at once).
3. Your cursor will turn into a crosshair.
4. Click and drag to select **just the anomaly card box** (the red-bordered card with the anomaly details).
5. Release the mouse button.
6. You'll hear a camera shutter sound, and a screenshot file will appear on your desktop (usually named something like `Screen Shot 2024-01-01 at 10.00.00 AM.png`).
7. **Rename and move the file:**
   - Right-click the screenshot file on your desktop.
   - Select "Rename" and change the name to `demo-anomaly-card.png`.
   - Drag it into the `screenshots` folder in your project (or move it there using Finder).

---

### Option B: Chrome DevTools Element Screenshot (Most Precise)

1. **Right-click** directly on the anomaly card (the box with the red "ANOMALY DETECTED" badge).
2. In the menu that appears, click **"Inspect"** (or "Inspect Element").
3. A panel will open at the bottom or side of your browser showing the HTML code.
4. In that panel, find the line that says something like `<div class="anomaly-card">` (it will be highlighted).
5. **Right-click** on that `<div class="anomaly-card">` line in the code panel.
6. In the menu, look for **"Capture node screenshot"** and click it.
7. The screenshot will automatically download to your Downloads folder.
8. **Move and rename the file:**
   - Find the downloaded screenshot (check your Downloads folder).
   - Rename it to `demo-anomaly-card.png`.
   - Move it to the `screenshots` folder in your project.

---

### Option C: Window Screenshot (Then Crop Later)

1. Make sure the anomaly card is visible in your browser window.
2. Press **Shift-Command-4**.
3. Press the **Spacebar** (your cursor will turn into a camera icon).
4. Move your mouse over the browser window and **click**.
5. The entire browser window will be captured as a screenshot.
6. The screenshot will appear on your desktop.
7. **Rename and move:**
   - Rename it to `demo-anomaly-card.png`.
   - Move it to the `screenshots` folder.
   - **Note:** You may want to crop this image later to show just the anomaly card, but it will work for now.

---

## Part 4: Verify Everything Worked

### Step 7: Check the Screenshot Location

Make sure the file exists at:
```
/Volumes/VIXinSSD/driftlock/screenshots/demo-anomaly-card.png
```

You can verify this by:
- Opening Finder and navigating to the `screenshots` folder in your project.
- Or running this command in terminal:
  ```bash
  ls -la screenshots/demo-anomaly-card.png
  ```
  (You should see file details, not an error.)

---

### Step 8: Check README Renders the Image (Optional)

1. Open `README.md` in your code editor.
2. Look for a preview mode (most editors have this).
3. Scroll to the bottom where it says "Demo Anomaly".
4. You should see your screenshot displayed there.

**If the image doesn't show:** Make sure the file is named exactly `demo-anomaly-card.png` and is in the `screenshots` folder.

---

## Part 5: Commit to Git (Optional)

Only do this if you want to save the screenshot to your git repository.

### Step 9: Add the File to Git
```bash
git add screenshots/demo-anomaly-card.png
```

### Step 10: Commit the File
```bash
git commit -m "add real demo anomaly screenshot"
```

### Step 11: Push to Remote (if desired)
```bash
git push
```

---

## Troubleshooting

**Problem:** Demo says "Found 0 anomalies" or a very low number (less than 10).

**Solution:** Run the verification script:
```bash
./verify-yc-ready.sh
```
This will check if everything is set up correctly.

---

**Problem:** Can't find the anomaly card in the browser.

**Solution:** 
- Make sure you ran the demo successfully (Step 4).
- Try scrolling down the page - anomalies are listed in order.
- Use Command-F to search for "ANOMALY" (without "DETECTED") to find all occurrences.

---

**Problem:** Screenshot file isn't showing up.

**Solution:**
- Check your Desktop folder (for Option A and C).
- Check your Downloads folder (for Option B).
- Make sure you completed the screenshot action (heard the shutter sound or saw the download).

---

## Quick Reference Checklist

- [ ] Navigated to project directory
- [ ] Cleaned up old files
- [ ] Built the demo program
- [ ] Ran the demo (saw "Found X anomalies...")
- [ ] Opened demo-output.html in browser
- [ ] Found first anomaly using Command-F search
- [ ] Captured screenshot using chosen method
- [ ] Moved screenshot to screenshots/demo-anomaly-card.png
- [ ] Verified file exists
- [ ] (Optional) Committed to git

---

*That's it! You should now have a screenshot of the first anomaly card saved in your project.*

