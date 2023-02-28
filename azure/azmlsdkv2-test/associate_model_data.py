import argparse
import json
import subprocess
from dataclasses import dataclass

from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class NameAndVersion:
    name: str
    version: str

    @classmethod
    def from_str(cls, x: str) -> "NameAndVersion":
        (name, version) = x.split(":")
        return NameAndVersion(name, version)


@dataclass
class Args:
    subscription: str
    resource_group: str
    workspace: str
    model: NameAndVersion
    data: NameAndVersion


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
    parser.add_argument("model", help="<model>:<version>")
    parser.add_argument("data", help="<data>:<version>")

    args = parser.parse_args()
    return Args(
        subscription=args.subscription,
        resource_group=args.resource_group,
        workspace=args.workspace,
        model=NameAndVersion.from_str(args.model),
        data=NameAndVersion.from_str(args.data),
    )


args = parse_args()

ml_client = MLClient(DefaultAzureCredential(), args.subscription, args.resource_group, args.workspace)

model = ml_client.models.get(args.model.name, version=args.model.version)
data = ml_client.data.get(args.data.name, version=args.data.version)

model.tags["sample_data"] = f"{data.name}:{data.version}"
ml_client.models.create_or_update(model)
