import argparse
import os
import sys
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    config: str
    job: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("job", help="job to inspect")

    args = parser.parse_args()
    return Args(
        config=args.config,
        job=args.job,
    )


def main():
    args = parse_args()

    ml_client = MLClient.from_config(DefaultAzureCredential(), file_name=args.config)
    job = ml_client.jobs.get(args.job)
    job.dump(sys.stdout)


if __name__ == "__main__":
    main()
