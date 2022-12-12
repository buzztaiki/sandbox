import argparse
from dataclasses import dataclass
from enum import Enum

from azureml.core import Dataset, Datastore, Workspace
from azureml.core.authentication import AzureCliAuthentication


class Action(Enum):
    ADD = 1
    DELETE = 2

    @classmethod
    def from_str(cls, x: str) -> "Action":
        if x == "add":
            return Action.ADD
        elif x == "delete":
            return Action.DELETE
        else:
            raise RuntimeError(f"unknown action: {x}")


@dataclass
class Args:
    subscription: str
    resource_group: str
    workspace: str
    datastore: str
    action: Action
    lst_file: str


def parse_args() -> Args:
    parser = argparse.ArgumentParser()
    parser.add_argument("-s", help="subscription")
    parser.add_argument("-w", help="ml workspace")
    parser.add_argument("-g", help="resource group")
    parser.add_argument("-d", help="datastore")
    parser.add_argument("action", help="add|delete")
    parser.add_argument("lst", help="dataset list")
    args = parser.parse_args()
    return Args(
        subscription=args.s,
        resource_group=args.g,
        workspace=args.w,
        datastore=args.d,
        action=Action.from_str(args.action),
        lst_file=args.lst,
    )


def main():
    args = parse_args()

    auth = AzureCliAuthentication()
    mlws = Workspace(
        subscription_id=args.subscription,
        resource_group=args.resource_group,
        workspace_name=args.workspace,
        auth=auth
    )

    with open(args.lst_file) as f:
        for line in f:
            dataset_name = line.rstrip()
            if args.action == Action.ADD:
                datastore = Datastore.get(mlws, args.datastore)
                dataset = Dataset.File.from_files(
                    path=(datastore, dataset_name),
                    # path=[(args.datastore, dataset_name)],
                    validate=False,
                )
                dataset.register(
                    workspace=mlws,
                    name=dataset_name,
                    description=f"Automatically registered from {args.datastore}/{dataset_name}",
                    create_new_version=True
                )
            elif args.action == Action.DELETE:
                dataset = Dataset.get_by_name(
                    workspace=mlws,
                    name=dataset_name,
                )
                if dataset is not None:
                    dataset.unregister_all_versions()
            else:
                assert False


if __name__ == "__main__":
    main()
