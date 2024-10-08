# from distutils.sysconfig import customize_compiler
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib.lines as mlines
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

def visualize_line_chart2(df, x_name, y_name, hue_cat_name, style_cat_name, path):
    # x_name = "numConditions"
    # y_name = "executionTime_mean"
    # hue_cat_name = "treeHeight"
    # style_cat_name = "mode"

    agg_df = df
    # agg_df = agg_df[
    #     (agg_df.step == step) &
    #     # (agg_df.name=="SingleProof") &
    #     # (agg_df.numConditions <= 5) &
    #     (agg_df.client=="nargo") &
    #     (agg_df.step == step)
    # ]

    x_values = agg_df[x_name].unique()
    hue_cat_values = agg_df[hue_cat_name].unique()
    style_cat_values = agg_df[style_cat_name].unique()

    palette = sns.color_palette("bright", len(hue_cat_values))

    marker_styles = ["o", "s", "D"]
    dash_styles = ["-", "--"]
    colors = palette

    fig, ax = plt.subplots()
    for im, hue_cat_value in enumerate(hue_cat_values):
        for il, style_cat_value in enumerate(style_cat_values):
            cur_df = agg_df[(agg_df[hue_cat_name] == hue_cat_value) & (agg_df[style_cat_name] == style_cat_value)]
            ax.plot(cur_df[x_name], cur_df[y_name], marker=marker_styles[im], linestyle=dash_styles[il], color=colors[im])

    legend_handles = []
    legend_handles.append(mlines.Line2D([0], [0], linestyle="none", marker="", label=get_title(hue_cat_name)))
    for im, heu_cat_value in enumerate(hue_cat_values):
        handle = mlines.Line2D([], [], color=colors[im], marker=marker_styles[im], label=heu_cat_value)
        legend_handles.append(handle)

    legend_handles.append(mlines.Line2D([0], [0], linestyle="none", marker="", label=get_title(style_cat_name)))
    for il, style_cat_value in enumerate(style_cat_values):
        handle = mlines.Line2D([], [], linestyle=dash_styles[il], color=colors[0], label=style_cat_value)
        legend_handles.append(handle)


    plt.legend(handles=legend_handles)
    plt.ylabel(get_title(y_name))
    plt.xlabel(get_title(x_name))
    plt.xticks(x_values)
    plt.grid(linestyle="--", axis="y", color="grey", linewidth=0.5)

    save_figure(fig, path)

    # if display:
    plt.show()

FONT_SIZE = 15
LABEL_SIZE = 15
LEGEND_SIZE = 15
MARKER_SIZE = 14

def visualize_line_chart(df, x_name, y_name, cat_name, path):
    x_values = df[x_name].unique()
    cat_values= df[cat_name].unique()

    logger.debug("x: {} - values: {}".format(x_name, x_values))
    logger.debug("cat: {} - values: {}".format(cat_name, cat_values))

    # sns.set_palette("pastel")
    custom_palette = sns.color_palette("bright", len(cat_values))
    sns.set_palette(custom_palette)
    # sns.palplot(custom_palette)
    figure = sns.lineplot(data=df, y=y_name, x=x_name, hue=cat_name, style=cat_name, palette=custom_palette, markers=True,
                          markersize=MARKER_SIZE).get_figure()

    plt.ylabel(get_title(y_name), fontsize=LABEL_SIZE)
    plt.xlabel(get_title(x_name), fontsize=LABEL_SIZE)
    plt.grid(linestyle="--", axis="y", color="grey", linewidth=0.5)
    plt.xticks(x_values, fontsize=FONT_SIZE)
    plt.yticks(fontsize=FONT_SIZE)
    plt.legend(title=get_title(cat_name), fontsize=LEGEND_SIZE, title_fontsize=FONT_SIZE)

    if path is not None:
        save_figure(figure, path)

    plt.show()

def visualize_bar_chart(df, x_name, y_name, cat_name, path):
    x_values = df[x_name].unique()
    y_values = df[y_name].unique()
    cat_values= df[cat_name].unique()

    logger.debug("x: {} - values: {}".format(x_name, x_values))
    logger.debug("y: {} - values: {}".format(x_name, y_values))
    logger.debug("cat: {} - values: {}".format(cat_name, cat_values))

    # sns.set_palette("pastel")
    custom_palette = sns.color_palette("bright", len(cat_values))
    sns.set_palette(custom_palette)
    # sns.palplot(custom_palette)
    figure = sns.barplot(data=df, y=y_name, x=x_name, hue=cat_name, palette=custom_palette).get_figure()

    plt.ylabel(get_title(y_name), fontsize=FONT_SIZE)
    plt.xlabel(get_title(x_name), fontsize=FONT_SIZE)
    plt.grid(linestyle="--", axis="y", color="grey", linewidth=0.5)
    plt.xticks(fontsize=FONT_SIZE)
    plt.legend(title=get_title(cat_name), fontsize=LEGEND_SIZE, title_fontsize=FONT_SIZE)
    plt.tight_layout()

    if path is not None:
        save_figure(figure, path)

    plt.show()

def save_figure(figure, path):
    if not os.path.exists(os.path.dirname(path)):
        os.makedirs(os.path.dirname(path))

    logger.info("saving figure to: {}".format(path))
    figure.savefig(path, bbox_inches='tight')

def get_title(name):
    name_dict = {
        "tps": "Throughput (Txs/second)",
        "avgLatency": "Average Latency (seconds)",
        "numOrgs": "Organizations",
        "numPeers": "Peers/Organization",
        "numVerifiers": "Verifiers",
        "contract": "Smart contract",
        "numTrans": "Token Transactions",
        "latency": "Average Latency (seconds)",
        "numTransTitle": "Token Transactions",
        "latencyM": "Average Latency (minutes)",
        "latency-infer": "Average Latency (seconds)",
        "latency-inferM": "Average Latency (minutes)",
        "short_contract": "Smart contract",
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

    visualize_line_chart(df, "numVerifiers", "tps", "short_contract", tps_figure_path)
    visualize_line_chart(df, "numVerifiers", "avgLatency", "short_contract", latency_figure_path)


def visualize_verification(df):
    partial_figure_path = os.path.join("..","..","exp_data","verification-partial.pdf")
    full_figure_path = os.path.join("..","..","exp_data","verification-full.pdf")
    logger.debug(df.columns)

    partial_df = df[df["contract"] == "SC_Verification#partial"]
    # visualize_bar_chart(partial_df, "numTransTitle", "latencyM", "numVerifiers", partial_figure_path)
    visualize_bar_chart(partial_df, "numVerifiers", "latency-inferM", "numTransTitle", partial_figure_path)

    full_df = df[df["contract"] == "SC_Verification#full"]
    # visualize_line_chart(full_df, "numVerifiers", "latency-inferM", "numTransTitle", full_figure_path)
    visualize_bar_chart(full_df, "numVerifiers", "latency-inferM", "numTransTitle", full_figure_path)

def visualize_collection(df):
    tps_collection_figure_path = os.path.join("..","..","exp_data","collection-tps.pdf")

    latency_collection_figure_path = os.path.join("..","..","exp_data","collection-latency.pdf")
    logger.debug(df.columns)

    # partial_df = df[df["contract"] == "SC_Verification#partial"]
    # visualize_bar_chart(partial_df, "numTransTitle", "latencyM", "numVerifiers", partial_figure_path)
    # visualize_bar_chart(df, "numVerifiers", "tps", "numVerifiers", collection_figure_path)
    visualize_line_chart(df, "numOrgs", "tps", "numVerifiers", tps_collection_figure_path)

    visualize_line_chart(df, "numOrgs", "avgLatency", "numVerifiers", latency_collection_figure_path)



def load_exp_data(exp_name):
    load_data_dict = {
        "caminit": "createCampaign.csv",
        "all": "all.csv",
        "collect": "collect.csv",
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
        "collect": visualize_collection,
    }

    visualize_fn_dict[exp_name](df)

def categorise(row):
    if row['numTrans'] == 115750:
        return "115,750 (1 week)"
    elif row["numTrans"] == 231500:
        return "231,500 (2 weeks)"
    elif row["numTrans"] == 463000:
        return "463,000 (4 weeks)"
    return "463,000 (4 weeks)"

def add_more_data(df):
    if "numTrans" in df.columns:
        df["numTransTitle"] = df.apply(lambda row: categorise(row), axis=1)

        df["latency-infer"] = df["latency"] / df["rawnumTrans"] * df["numTrans"]
        df["latencyM"] = df["latency"] / 60
        df["latency-inferM"] = df["latency-infer"] / 60

    # df.loc[df["numTrans"] == 997260, "numTransTitle"] = "997,260 (1 week)"
    # df.loc[df["numTrans"] == 1994520, "numTransTitle"] = "1,994,520 (2 weeks)"
    # df.loc[df["numTrans"] == 3989040, "numTransTitle"] = "3,989,040 (4 weeks)"

    # df["numTransTitle"] = df["numTransTitle"].astype(str)
    def get_short_contract_name(name):
        name_map = {
            "SC_Witness_Token_Generation": "SC_WTGen",
            "SC_Campaign_Token_Generation": "SC_CTGen",
            "SC_Initialization": "SC_Init",
            "SC_Token_Collection": "SC_TCol",
            "SC_Verification#full": "SC_VerFull",
            "SC_Verification#partial": "SC_VerPar",
        }

        return name_map[name]

    if "contract" in df.columns:
        df["short_contract"] = df["contract"].apply(get_short_contract_name)

    pass

def main(args):
    exp_name = args["exp"]

    df = load_exp_data(exp_name)

    add_more_data(df)

    logger.debug(df)
    df.to_csv("temp.csv")
    visualize(exp_name, df)

if __name__ == "__main__":
    args = setup_arguments(add_arguments)
    setup_console_logging(args)
    main(args)
