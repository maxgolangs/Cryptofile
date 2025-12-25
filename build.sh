#!/bin/bash

echo "🔐 Сборка бинарников CryptoFile by @MaxGolang..."

# Создаем папку dist
mkdir -p dist

echo "🧹 Очистка..."
rm -f dist/cryptofile-linux-amd64 dist/cryptofile-windows-amd64.exe

# Сборка для Linux (amd64)
echo ""
echo "🔨 Сборка для Linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o dist/cryptofile-linux-amd64 ./cmd/cryptor
chmod +x dist/cryptofile-linux-amd64
echo "✓ Linux бинарник создан: dist/cryptofile-linux-amd64"

# Сборка для Windows (amd64)
echo ""
echo "🔨 Сборка для Windows/amd64..."

# Проверяем наличие MinGW для кросс-компиляции
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "⚠️  MinGW не найден. Установка требуется для сборки Windows версии..."
    echo "💡 Для установки выполните: sudo apt install gcc-mingw-w64"
    echo "⚠️  Пропускаем сборку Windows версии"
else
    echo "✓ MinGW найден, используем CGO для сборки..."
    CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
        -o dist/cryptofile-windows-amd64.exe ./cmd/cryptor 2>&1
    
    if [ -f dist/cryptofile-windows-amd64.exe ]; then
        echo "✓ Windows бинарник создан: dist/cryptofile-windows-amd64.exe"
    else
        echo "⚠️  Ошибка при сборке Windows бинарника"
    fi
fi

echo ""
echo "✅ Готово! Все бинарники собраны в папке dist/:"
ls -lh dist/

