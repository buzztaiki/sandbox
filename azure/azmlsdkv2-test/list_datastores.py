import argparse
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    config: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    args = parser.parse_args()
    return Args(config=args.config)


def main():
    args = parse_args()

    ml_client = MLClient(DefaultAzureCredential(), file_name=args.config)
    for x in ml_client.datastores.list():
        print(x.name, x.type)


if __name__ == "__main__":
    main()
