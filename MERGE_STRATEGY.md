# Branch Merge Strategy

## Current Situation
- **main**: Working CLI demo + verification script (as per FINAL-STATUS.md)
- **cleanup-and-transition**: Removed demo, prepared for SaaS transition 
- **landing-page-professional-improvements**: Landing page enhancements
- **saas-launch**: AI continuation prompt

## Problem
The cleanup-and-transition branch violated the "Golden Invariant" by removing the working demo.

## Solution: Preserve Demo + Add SaaS Features

### Step 1: Merge landing page improvements to main
These are safe improvements that don't break the demo.

### Step 2: Selectively merge SaaS features from cleanup branch
- Keep: Firebase config, Cloud Build setup, deployment docs
- Keep: Landing page SaaS enhancements  
- Skip: Removal of cbad-core, cmd/demo, verify script
- Keep: New SaaS architecture docs

### Step 3: Create hybrid approach
- Maintain CLI demo for verification/partners
- Add SaaS platform alongside (not instead of)
- Both approaches can coexist

## Files to Merge

### From landing-page-professional-improvements:
- All landing-page/ improvements (SEO, accessibility, form validation)

### From cleanup-and-transition (selective):
- .firebaserc (Firebase config)
- .gcloudignore (deployment)  
- cloudbuild.yaml (Cloud Build)
- deploy.sh (deployment script)
- docs/deployment/ (new deployment docs)
- docs/STREAMING.md (new docs)
- Firebase/Cloud Run setup docs
- Landing page routing enhancements

### Keep from main:
- All cbad-core/ (Rust implementation)
- All cmd/demo/ (CLI demo)
- All collector-processor/ (Go HTTP API)
- verify-yc-ready.sh (verification script)
- Makefile (build system)
- test-data/ (demo data)

## Result
A hybrid repository that:
1. ✅ Maintains working CLI demo (satisfies AGENTS.md)
2. ✅ Adds SaaS deployment capability
3. ✅ Has professional landing page
4. ✅ Can be deployed to Cloud Run/Firebase
5. ✅ Preserves all roadmap capabilities