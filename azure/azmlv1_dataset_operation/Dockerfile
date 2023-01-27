FROM ubuntu:20.04
RUN apt-get update && apt-get install -y --no-install-recommends python3 python3-pip curl
RUN pip install poetry
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

WORKDIR /app
COPY poetry.lock pyproject.toml /app
RUN poetry install
