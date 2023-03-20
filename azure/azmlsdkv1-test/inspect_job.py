import argparse
import pprint
from dataclasses import dataclass

from azureml.core import Dataset, Run, Workspace


@dataclass
class Args:
    config: str
    job: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("job", help="job to inspect")

    args = parser.parse_args()
    return Args(config=args.config, job=args.job)


def main():
    args = parse_args()
    ws = Workspace.from_config(_file_name=args.config)
    run = Run.get(ws, args.job)

    print("JOB:")
    print(run)
    print()

    details = run.get_details()
    print("JOB_DETAILS:")
    pprint.pp(details)
    print()

    print("INPUT_DATASETS:")
    for x in details["inputDatasets"]:
        ds: Dataset = x["dataset"]
        print(ds.name, ds.id)


if __name__ == '__main__':
    main()
