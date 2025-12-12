FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

RUN echo "test-results.xml" >> .dockerignore
RUN npm run build

EXPOSE 3000

# Install wait-port globally
RUN npm install -g wait-port

# Wait for backend BEFORE starting Next.js
CMD ["sh", "-c", "wait-port backend:10000 && npm start"]
