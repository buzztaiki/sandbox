from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential
import mlflow


def main():
    ml_client = MLClient.from_config(DefaultAzureCredential())
    ws = ml_client.workspaces.get(ml_client.workspace_name)
    mlflow.set_tracking_uri(ws.mlflow_tracking_uri)

    for run in mlflow.search_runs(
            search_all_experiments=True,
            output_format="list"
    ):
        print(run)


if __name__ == "__main__":
    main()
