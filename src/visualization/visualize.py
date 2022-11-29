from distutils.sysconfig import customize_compiler
import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
import logging
import argparse
import os

logging.getLogger("matplotlib").setLevel(logging.WARNING)
logger = logging.getLogger(__name__)

def str2log_mode(value):
    if value is None:
        return None

    if value in ["d", "debug", "10"]:
        log_mode = logging.DEBUG
    elif value in ["i", "info", "20"]:
        log_mode = logging.INFO
    elif value in ["w", "warning", "30"]:
        log_mode = logging.WARNING
    else:
        raise argparse.ArgumentTypeError("Unsupported log mode type: {}".format(value))

    return log_mode

def setup_arguments(add_arguments_fn):
    parser = argparse.ArgumentParser(description="Process some integers.")

    parser.add_argument("--log", type=str2log_mode, default=logging.INFO)
    add_arguments_fn(parser)

    args, _ = parser.parse_known_args()

    params = {}
    for arg in vars(args):
        params[arg] = getattr(args, arg)

    # os.environ[ASSERTION_VARIABLE] = params["assert"]

    return params

def setup_console_logging(args):
    level = args["log"]

    logger = logging.getLogger("")
    logger.setLevel(level)

    formatter = logging.Formatter(
        "%(name)-12s[%(lineno)d]: %(funcName)s %(levelname)-8s %(message)s "
    )

    console_handler = logging.StreamHandler()
    console_handler.setLevel(level)
    console_handler.setFormatter(formatter)

    logger.addHandler(console_handler)

def add_arguments(parser):
    parser.add_argument("--exp")

def visualize_line_chart(df, x_name, y_name, cat_name, path):
    x_values = df[x_name].unique()
    cat_values= df[cat_name].unique()

    logger.debug("x: {} - values: {}".format(x_name, x_values))
    logger.debug("cat: {} - values: {}".format(cat_name, cat_values))

    # sns.set_palette("pastel")
    custom_palette = sns.color_palette("bright", len(cat_values))
    sns.set_palette(custom_palette)
    # sns.palplot(custom_palette)
    figure = sns.lineplot(data=df, y=y_name, x=x_name, hue=cat_name, style=cat_name, palette=custom_palette, markers=True).get_figure()

    plt.ylabel(get_title(y_name))
    plt.xlabel(get_title(x_name))
    plt.grid(linestyle="--", axis="y", color="grey", linewidth=0.5)
    plt.xticks(x_values)
    plt.legend(title=get_title(cat_name))

    if path is not None:
        save_figure(figure, path)

    plt.show()

def save_figure(figure, path):
    if not os.path.exists(os.path.dirname(path)):
        os.makedirs(os.path.dirname(path))

    logger.info("saving figure to: {}".format(path))
    figure.savefig(path)

def get_title(name):
    name_dict = {
        "tps": "Throughput (Txs/second)",
        "avgLatency": "Average Latency (seconds)",
        "numOrgs": "# of Organizations",
        "numPeers": "# of Peers per Organization",
        "numVerifiers": "# of Verifiers",
        "contract": "Smart contract",
        "numTrans": "# of Token Transactions"
    }

    return name_dict[name]

def visualize_campaign_init(df):
    tps_figure_path = os.path.join("..","..","exp_data","caminit-tps.pdf")
    latency_figure_path = os.path.join("..","..","exp_data","caminit-latency.pdf")
    logger.debug(df.columns)

    visualize_line_chart(df, "numOrgs", "tps", "numPeers", tps_figure_path)
    visualize_line_chart(df, "numOrgs", "avgLatency", "numPeers", latency_figure_path)


def visualize_all(df):
    tps_figure_path = os.path.join("..","..","exp_data","all-tps.pdf")
    latency_figure_path = os.path.join("..","..","exp_data","all-latency.pdf")
    logger.debug(df.columns)

    visualize_line_chart(df, "numVerifiers", "tps", "contract", tps_figure_path)
    visualize_line_chart(df, "numVerifiers", "avgLatency", "contract", latency_figure_path)


def visualize_verification(df):
    partial_figure_path = os.path.join("..","..","exp_data","verification-partial.pdf")
    full_figure_path = os.path.join("..","..","exp_data","verification-full.pdf")
    logger.debug(df.columns)

    partial_df = df[df["contract"] == "SC_Verification#partial"]
    visualize_line_chart(partial_df, "numVerifiers", "avgLatency", "numTrans", partial_figure_path)

    full_df = df[df["contract"] == "SC_Verification#partial"]
    visualize_line_chart(full_df, "numVerifiers", "avgLatency", "numTrans", full_figure_path)

def load_exp_data(exp_name):
    load_data_dict = {
        "caminit": "createCampaign.csv",
        "all": "all.csv",
        "ver": "verification.csv",
    }

    path = os.path.join("..", "..", "exp_data", load_data_dict[exp_name])
    logger.debug(path)

    df = pd.read_csv(path)
    return df

def visualize(exp_name, df):
    visualize_fn_dict = {
        "caminit": visualize_campaign_init,
        "all": visualize_all,
        "ver": visualize_verification,
    }

    visualize_fn_dict[exp_name](df)

def main(args):
    exp_name = args["exp"]

    df = load_exp_data(exp_name)

    logger.debug(df)
    visualize(exp_name, df)

if __name__ == "__main__":
    args = setup_arguments(add_arguments)
    setup_console_logging(args)
    main(args)
