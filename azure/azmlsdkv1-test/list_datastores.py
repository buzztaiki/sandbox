import argparse
from dataclasses import dataclass

from azureml.core import Datastore, Workspace


@dataclass
class Args:
    config: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")

    args = parser.parse_args()
    return Args(
        config=args.config,
    )


def main():
    args = parse_args()
    ws = Workspace.from_config(_file_name=args.config)
    for x in ws.datastores:
        # ds: Union[Datastore, None] = ws.datastores.get(x)
        ds = Datastore.get(ws, x)
        if ds is not None:
            print(ds.name)


main()
