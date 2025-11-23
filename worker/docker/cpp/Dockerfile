FROM gcc:12.2.0

# Create workspace folder
WORKDIR /app

# Copy runner script
COPY run.sh /app/run.sh
RUN chmod +x /app/run.sh

# This container will compile and run C++ code provided via env variable CODE
ENTRYPOINT ["/app/run.sh"]
