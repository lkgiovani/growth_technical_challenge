#!/bin/bash

set -e

echo "🔍 Checking Node.js installation..."
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ and try again."
    exit 1
fi

echo "📦 Installing dependencies..."
npm install

echo "🔨 Compiling TypeSpec to OpenAPI..."
npm run compile

echo "✅ Compilation complete!"
echo ""
echo "Generated files:"
echo "  - docs/schema/openapi.yaml"
echo ""
echo "TypeSpec source files are in:"
echo "  - docs/src/"
echo ""
echo "To view the documentation, start your Go API and visit:"
echo "  - http://localhost:8080/docs/redoc"
echo "  - http://localhost:8080/docs/swagger"
echo "  - http://localhost:8080/docs/scalar"

