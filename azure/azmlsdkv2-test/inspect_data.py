import argparse
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    config: str
    data: str
    version: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-v", "--version", help="model version (default: latest)", default="latest")
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("data", help="data to inspect")

    args = parser.parse_args()
    return Args(
        config=args.config,
        data=args.data,
        version=args.version,
    )


def main():
    args = parse_args()

    ml_client = MLClient.from_config(DefaultAzureCredential(), file_name=args.config)
    if args.version == "latest":
        data = ml_client.data.get(args.data, label="latest")
    else:
        data = ml_client.data.get(args.data, version=args.version)
    data.print_as_yaml = True
    print(data)


if __name__ == "__main__":
    main()
