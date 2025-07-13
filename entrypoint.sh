#!/bin/sh

echo "⏳ Waiting for MySQL to be ready..."
/wait-for-it.sh db:3306 --timeout=60 --strict -- echo "✅ MySQL is up"

echo "🚀 Starting Go app..."
./main