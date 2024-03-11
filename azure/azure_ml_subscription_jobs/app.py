import itertools
import json
from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any, Generator

from azure.ai.ml import MLClient
from azure.ai.ml.entities import Job
from azure.core.credentials import TokenCredential
from azure.identity import DefaultAzureCredential
from azure.mgmt.resourcegraph import ResourceGraphClient
from azure.mgmt.resourcegraph.models import QueryRequest


@dataclass
class WorkspaceInfo:
    subscription_id: str
    resource_group: str
    name: str


def workspaces(cred: TokenCredential) -> Generator[WorkspaceInfo, None, None]:
    graph_client = ResourceGraphClient(cred)
    query = 'resources | where type == "microsoft.machinelearningservices/workspaces"'
    data: Any = graph_client.resources(QueryRequest(query=query)).data
    for res in data:
        yield WorkspaceInfo(res["subscriptionId"], res["resourceGroup"], res["name"])


def job_start_dt(job: Job) -> datetime | None:
    match job.type:
        case "command":
            if dtstr := job.properties.get("StartTimeUtc"):
                return datetime.strptime(dtstr, "%Y-%m-%d %H:%M:%S").replace(
                    tzinfo=timezone.utc
                )
        case "pipeline":
            x = json.loads(job.properties["azureml.pipelines.stages"])
            if dtstr := x["Execution"].get("StartTime"):
                return datetime.fromisoformat(dtstr)
        case "base":
            # TODO: base の場合に日時が入ってる属性がなさげ
            pass
        case _:
            pass
    return None


def job_end_dt(job: Job) -> datetime | None:
    match job.type:
        case "command":
            if dtstr := job.properties.get("EndTimeUtc"):
                return datetime.strptime(dtstr, "%Y-%m-%d %H:%M:%S").replace(
                    tzinfo=timezone.utc
                )
        case "pipeline":
            x = json.loads(job.properties["azureml.pipelines.stages"])
            if dtstr := x["Execution"].get("EndTime"):
                return datetime.fromisoformat(dtstr)
        case "base":
            # TODO: base の場合に日時が入ってる属性がなさげ
            pass
        case _:
            pass
    return None


def iso8601(dt: datetime | None) -> str | None:
    if dt is not None:
        return dt.isoformat()
    return None


@dataclass
class JobWrapper:
    job: Job
    parent_job_name: str | None


def jobs(
    client: MLClient, parent_job_name: str | None = None
) -> Generator[JobWrapper, None, None]:
    for job in client.jobs.list(parent_job_name=parent_job_name):
        if job.type == "pipeline":
            yield JobWrapper(job, parent_job_name)
            for child in jobs(client, job.name):
                yield child
        else:
            yield JobWrapper(job, parent_job_name)


def main():
    cred = DefaultAzureCredential()
    for wsinfo in workspaces(cred):
        client = MLClient(
            cred,
            subscription_id=wsinfo.subscription_id,
            resource_group_name=wsinfo.resource_group,
            workspace_name=wsinfo.name,
        )
        client.jobs.list()
        for job in itertools.islice(jobs(client), 10):
            cols: list[str | None] = [
                client.workspace_name,
                job.job.name,
                job.job.display_name,
                job.job.type,
                job.parent_job_name,
                job.job.compute,
                iso8601(job_start_dt(job.job)),
                iso8601(job_end_dt(job.job)),
            ]
            print("\t".join([str(x) if x is not None else "" for x in cols]))
            # buf = StringIO()
            # job.job.dump(dest=buf)
            # print(buf.getvalue())
            # print()


if __name__ == "__main__":
    main()
