FROM ubuntu:latest
RUN apt update -y && apt install -y nodejs npm curl
RUN npm install -g pnpm
WORKDIR /frontend
COPY package.json ./frontend/pnpm-lock.yaml ./
RUN pnpm install --force
COPY ./frontend .
EXPOSE 3000
RUN pnpm run build
CMD ["sh", "-c", "pnpm run dev & tail -f /dev/null"]
