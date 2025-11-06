#!/bin/bash

# script para preparar entorno de test aislado
# no afecta datos de produccion

echo "preparando entorno de test..."
echo "=============================="

# crear directorio de test
mkdir -p backend/data-test/prompts

# copiar faqs a test
if [ -f "backend/data/faqs.csv" ]; then
    cp backend/data/faqs.csv backend/data-test/faqs.csv
    echo "✓ faqs copiadas a data-test"
else
    echo "⚠ no se encontro backend/data/faqs.csv"
fi

# crear archivos vacios para sessions y leads
echo "{}" > backend/data-test/sessions.json
echo "{}" > backend/data-test/leads.json
echo "✓ sessions.json y leads.json inicializados"

# crear .env.test
cat > backend/.env.test << EOF
GEMINI_API_KEY=$(grep GEMINI_API_KEY backend/.env | cut -d '=' -f2)
GEMINI_MODEL=gemini-2.5-flash
PORT=3001
BOB_API_BASE_URL=https://apiv3.somosbob.com/v3
CORS_ORIGINS=http://localhost:5173,http://localhost:3001
FRONTEND_URL=http://localhost:5173
DATA_DIR=data-test
EOF

echo "✓ .env.test creado (puerto 3001, data-test)"
echo ""
echo "=============================="
echo "entorno de test listo"
echo "=============================="
echo ""
echo "para iniciar servidor de test:"
echo "  cd backend"
echo "  cp .env.test .env"
echo "  go run cmd/server/main.go"
echo ""
echo "o directamente:"
echo "  cd backend && export \$(cat .env.test | xargs) && go run cmd/server/main.go"
echo ""
echo "el servidor de test:"
echo "  - corre en puerto 3001"
echo "  - usa data-test/ (no afecta data/)"
echo "  - sessions y leads empiezan vacios"
echo "  - faqs son las mismas"
