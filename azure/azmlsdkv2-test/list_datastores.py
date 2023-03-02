import argparse
import json
import subprocess
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    subscription: str
    resource_group: str
    workspace: str


def parse_args():
    subscription = json.loads(subprocess.run(["az", "account", "show"], capture_output=True, check=True).stdout)["id"]
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-s",
        "--subscription",
        help=f"subscription (default: {subscription})",
        default=subscription,
    )
    parser.add_argument("-g", "--resource-group", help="resource group", required=True)
    parser.add_argument("-w", "--workspace", help="ml workspace", required=True)

    args = parser.parse_args()
    return Args(
        subscription=args.subscription,
        resource_group=args.resource_group,
        workspace=args.workspace,
    )


def main():
    args = parse_args()

    ml_client = MLClient(DefaultAzureCredential(), args.subscription, args.resource_group, args.workspace)
    for x in ml_client.datastores.list():
        print(x.name, x.type)


if __name__ == "__main__":
    main()
