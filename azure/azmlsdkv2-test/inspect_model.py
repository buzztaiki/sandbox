import argparse
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    config: str
    model: str
    version: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("-v", "--version", help="model version (default: latest)", default="latest")
    parser.add_argument("model", help="model to inspect")

    args = parser.parse_args()
    return Args(
        config=args.config,
        model=args.model,
        version=args.version,
    )


def main():
    args = parse_args()
    ml_client = MLClient.from_config(DefaultAzureCredential(), file_name=args.config)

    if args.version == "latest":
        model = ml_client.models.get(args.model, label="latest")
    else:
        model = ml_client.models.get(args.model, version=args.version)

    model.print_as_yaml = True
    print(model)


if __name__ == "__main__":
    main()
