# https://github.com/mlflow/mlflow/tree/master/examples/sklearn_elasticnet_wine
# https://mlflow.org/docs/latest/tutorials-and-examples/tutorial.html

import argparse
import json
import logging
import warnings
from dataclasses import dataclass

import mlflow
import mlflow.sklearn
import numpy as np
import pandas as pd
from sklearn.linear_model import ElasticNet
from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score
from sklearn.model_selection import train_test_split

logging.basicConfig(level=logging.WARN)
logger = logging.getLogger(__name__)


def eval_metrics(actual, pred):
    rmse = np.sqrt(mean_squared_error(actual, pred))
    mae = mean_absolute_error(actual, pred)
    r2 = r2_score(actual, pred)
    return rmse, mae, r2


@dataclass
class Args:
    experiment: str
    alpha: float
    l1_ratio: float


def parse_args() -> Args:
    parser = argparse.ArgumentParser()
    parser.add_argument("-e", "--experiment", help="experiment name", required=True)
    parser.add_argument("--alpha", type=float, default=0.5)
    parser.add_argument("--l1_ratio", type=float, default=0.5)

    args = parser.parse_args()
    return Args(experiment=args.experiment,
                alpha=args.alpha,
                l1_ratio=args.l1_ratio)


def main():
    args = parse_args()
    mlflow.set_experiment(args.experiment)

    warnings.filterwarnings("ignore")
    np.random.seed(40)

    # Read the wine-quality csv file from the URL
    csv_url = "https://raw.githubusercontent.com/mlflow/mlflow/master/tests/datasets/winequality-red.csv"
    try:
        data = pd.read_csv(csv_url, sep=";")
    except Exception as e:
        logger.exception(
            "Unable to download training & test CSV, check your internet connection. Error: %s",
            e,
        )

    # Split the data into training and test sets. (0.75, 0.25) split.
    train, test = train_test_split(data)

    # The predicted column is "quality" which is a scalar from [3, 9]
    train_x = train.drop(["quality"], axis=1)
    test_x = test.drop(["quality"], axis=1)
    train_y = train[["quality"]]
    test_y = test[["quality"]]

    alpha = args.alpha
    l1_ratio = args.l1_ratio

    with mlflow.start_run():
        mlflow.set_tag("data", json.dumps({"csv_url": csv_url}))

        lr = ElasticNet(alpha=alpha, l1_ratio=l1_ratio, random_state=42)
        lr.fit(train_x, train_y)

        predicted_qualities = lr.predict(test_x)

        (rmse, mae, r2) = eval_metrics(test_y, predicted_qualities)

        print("Elasticnet model (alpha={:f}, l1_ratio={:f}):".format(alpha, l1_ratio))
        print("  RMSE: %s" % rmse)
        print("  MAE: %s" % mae)
        print("  R2: %s" % r2)

        mlflow.log_param("alpha", alpha)
        mlflow.log_param("l1_ratio", l1_ratio)
        mlflow.log_metric("rmse", rmse)
        mlflow.log_metric("r2", r2)
        mlflow.log_metric("mae", mae)

        # mlflow.sklearn.save_model(lr, path="trained_model")
        mlflow.sklearn.log_model(lr, "model")
        # mlflow.sklearn.log_model(
        #     lr, "model", registered_model_name="ElasticnetWineModel"
        # )


if __name__ == "__main__":
    main()
