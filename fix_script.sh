#!/bin/bash
# File to fix specific issues in the attendance_check_in.go file

FILE_PATH="/Users/tinhp/Workspace/ITC/PROJECTS/gomcms/server/service/checkins/attendance_check_in.go"

# Use sed to fix specific indentation issues
sed -i '' 's/resultList := checkConditionsParallel(gc.agp, gc.conditions, req, ip)\t\t\t\t\/\/ Xử lý kết quả/resultList := checkConditionsParallel(gc.agp, gc.conditions, req, ip) \/\/ Xử lý kết quả/g' "$FILE_PATH"
sed -i '' 's/\t\t\t\t\t}\t\t\t\t\tif result.Pass {/\t\t\t\t\t}\n\n\t\t\t\t\tif result.Pass {/g' "$FILE_PATH"

echo "File fixed successfully!"
