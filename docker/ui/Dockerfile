# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files and install deps
COPY ui/package.json ui/yarn.lock ./

# Copy only source code
COPY ui/.babel-plugin-macrosrc.js ./
COPY ui/.dockerignore ./
COPY ui/.editorconfig ./
COPY ui/.env.production ./
COPY ui/.eslintignore ./
COPY ui/.eslintrc.js ./
COPY ui/.gitattributes ./
COPY ui/.gitignore ./
COPY ui/.npmrc ./
COPY ui/.nvmrc ./
COPY ui/.pre-commit-config.yaml ./
COPY ui/.prettierignore ./
COPY ui/.prettierrc ./
COPY ui/.stylelintrc ./
COPY ui/amplify.yml ./
COPY ui/commitlint.config.js ./
COPY ui/craco.config.js ./
COPY ui/electron-builder.yml ./
COPY ui/package.json ./
COPY ui/postcss.config.js ./
COPY ui/template.json ./
COPY ui/tsconfig.json ./
COPY ui/tsconfig.paths.json ./
COPY ui/yarn.lock ./

# RUN yarn install --frozen-lockfile
RUN yarn install --frozen-lockfile --cache-folder .yarn-cache && \
    rm -rf .yarn-cache

# Copy directories
COPY ui/public/ ./public
COPY ui/src/ ./src


# Build the React app
RUN yarn build:docker

# # Final stage - Nginx
FROM nginx:alpine

# Create user
RUN addgroup -g 1000 rapida-app && \
    adduser -D -u 1000 -G rapida-app rapida-app

# Copy nginx config
COPY docker/nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built app
COPY --from=builder /app/build /usr/share/nginx/html

# Ensure ownership
RUN chown -R rapida-app:rapida-app /usr/share/nginx/html /var/cache/nginx /var/log/nginx
RUN mkdir -p /run && chown -R rapida-app:rapida-app /run


# Switch to user
USER rapida-app

# Expose port
EXPOSE 3000

# Run nginx
CMD ["nginx", "-g", "daemon off;"]