import argparse
import os
from dataclasses import dataclass
from random import randint, random

import mlflow


@dataclass
class Args:
    experiment: str


def parse_args() -> Args:
    parser = argparse.ArgumentParser()
    parser.add_argument("-e", "--experiment", help="experiment name", required=True)

    args = parser.parse_args()
    return Args(experiment=args.experiment)


def main():
    args = parse_args()
    mlflow.set_experiment(args.experiment)
    # Log a parameter (key-value pair)
    mlflow.log_param("param1", randint(0, 100))

    # Log a metric; metrics can be updated throughout the run
    mlflow.log_metric("foo", random())
    mlflow.log_metric("foo", random() + 1)
    mlflow.log_metric("foo", random() + 2)

    # Log an artifact (output file)
    if not os.path.exists("outputs"):
        os.makedirs("outputs")
    with open("outputs/test.txt", "w") as f:
        f.write("hello world!")
    mlflow.log_artifacts("outputs")


if __name__ == "__main__":
    main()
