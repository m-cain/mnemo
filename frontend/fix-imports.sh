#!/bin/bash

# Fix @/lib/utils imports in UI components
find frontend/src/components/ui -type f -name "*.tsx" -exec sed -i '' 's/from "@\/lib\/utils"/from "..\/..\/lib\/utils"/g' {} \;

# Fix other @/ imports in UI components
find frontend/src/components/ui -type f -name "*.tsx" -exec sed -i '' 's/from "@\/components\/ui\/\([^"]*\)"/from ".\1"/g' {} \;

echo "Fixed imports in UI components"
