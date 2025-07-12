import argparse
import pprint
from dataclasses import dataclass

import mlflow
from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential


@dataclass
class Args:
    config: str
    run_name: str


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config json", default="config.json")
    parser.add_argument("run_name", help="run to inspect")

    args = parser.parse_args()
    return Args(
        config=args.config,
        run_name=args.run_name,
    )


def main():
    args = parse_args()
    ml_client = MLClient.from_config(DefaultAzureCredential(), file_name=args.config)
    ws = ml_client.workspaces.get(ml_client.workspace_name)

    assert ws is not None
    assert ws.mlflow_tracking_uri is not None

    # Note: we must use set_tracking_uri to set tracking and registry uri with azureml-mlflow extension
    mlflow.set_tracking_uri(ws.mlflow_tracking_uri)
    mlflow_client = mlflow.MlflowClient()
    run = mlflow_client.get_run(args.run_name)
    print("RUN:")
    pprint.pp(run.to_dictionary())
    print()

    print("ARTIFACTS:")
    run_info = run.info
    for artifact in mlflow_client.list_artifacts(run_info.run_id):
        print(artifact)
    print()

    print("MODELS:")
    for model in mlflow_client.search_model_versions(f"run_id='{run_info.run_id}'"):
        pprint.pp(model)
    print()


if __name__ == "__main__":
    main()
