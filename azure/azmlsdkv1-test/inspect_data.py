import argparse
import pprint
from dataclasses import dataclass

from azureml.core import Dataset, Workspace


@dataclass
class Args:
    config: str
    data: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("data", help="data to inspect")

    args = parser.parse_args()
    return Args(config=args.config, data=args.data)


def main():
    args = parse_args()
    ws = Workspace.from_config(_file_name=args.config)
    dataset = Dataset.get_by_name(ws, args.data)
    assert dataset is not None
    print(dataset)


if __name__ == '__main__':
    main()
