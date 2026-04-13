import pandas as pd
import os

file_keep_dict = {
    "occurrence.csv": [
        "occurrence_no",
        "collection_no",
        "accepted_name",
        "accepted_rank",
        "early_interval",
        "late_interval",
        "max_ma",
        "min_ma",
    ],
    "collection.csv": [
        "collection_no",
        "formation",
        "lng",
        "lat",
        "collection_name",
        "n_occs",
        "early_interval",
        "late_interval",
        "max_ma",
        "min_ma",
    ],
    "taxa.csv": [
        "taxon_no",
        "accepted_name",
        "accepted_rank",
        "parent_no",
        "is_extant",
        "n_occs",
    ],
}


def clean_data(file_name, keep_cols):

    path = "data/hwdata/" + file_name
    df = pd.read_csv(path)
    df_clean = df[keep_cols]

    # df_clean = df_clean.dropna(subset=["accepted_name", "max_ma", "min_ma"])

    df_clean = df_clean.drop_duplicates()

    os.makedirs("data/cleaned_hwdata", exist_ok=True)

    df_clean.to_csv("data/cleaned_hwdata/" + file_name, index=False)

    print("清洗完成")


for filename in file_keep_dict.keys():
    clean_data(filename, file_keep_dict[filename])
