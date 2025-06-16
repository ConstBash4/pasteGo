<div align="center">
<img src="./frontend/src/assets/logo_small.svg" width="256" alt=""/>
  <h1> PasteGo 📋</h1>
  <p> A simple text-sharing platform built with Go and SvelteKit </p>
</div>

### 🐳 Docker
```bash
git clone https://github.com/ConstBash4/pasteGo.git
cd pasteGo
#change SECRET_KEY in .env.example and rename it
mv .env.example .env
#choose standalone or compose
docker build --tag pastego .
docker run -d --restart=unless-stopped -p 10015:10015 -v ./data:/app/data --name pastego pastego:latest
#OR
docker compose up -d
