#!/bin/bash
# Fix the tiered storage file by replacing problematic lines

# Create a temporary file to hold the modified content
awk '
NR == 118 { 
  print "\t// Use filtered directly (no sorting since AnomalyFilter doesn't have SortField/SortOrder)"
  next
}
NR == 119 { 
  print "\tsorted := filtered"
  next
}
{ print }
' /Volumes/VIXinSSD/driftlock/api-server/internal/storage/tiered.go > /tmp/tiered_fixed.go

# Replace the original file
mv /tmp/tiered_fixed.go /Volumes/VIXinSSD/driftlock/api-server/internal/storage/tiered.go