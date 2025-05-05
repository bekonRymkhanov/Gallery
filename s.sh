echo "Fixing Firefox HTTPS redirection issue..."

# Create script to update index.html in the frontend container
cat > update-index.sh << 'EOL'
#!/bin/sh
INDEX_PATH="/usr/share/nginx/html/index.html"

# Backup original index.html
cp $INDEX_PATH ${INDEX_PATH}.bak

# Add script to prevent HTTPS redirects to index.html
sed -i 's/<head>/<head>\n  <!-- Prevent HTTPS redirects -->\n  <meta http-equiv="Content-Security-Policy" content="upgrade-insecure-requests '\''none'\''">/' $INDEX_PATH
sed -i 's/<head>/<head>\n  <script>\n    if (window.location.protocol === "https:") {\n      window.location.protocol = "http:";\n    }\n  <\/script>/' $INDEX_PATH

echo "Updated index.html to prevent HTTPS redirects"
EOL

# Make the script executable
chmod +x update-index.sh

# Copy the script to the nginx container
docker cp update-index.sh nginx:/tmp/update-index.sh

# Run the script in the nginx container
docker exec nginx sh /tmp/update-index.sh

# Clean up the temporary script
rm update-index.sh

echo "Firefox HTTPS redirection fix has been applied."
echo "Please clear your browser cache and try accessing http://localhost/ again."