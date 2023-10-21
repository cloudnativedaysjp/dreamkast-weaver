#!/usr/bin/env python3
from typing import Final, Literal, Optional

import fire
import pandas as pd
import matplotlib.pyplot as plt

from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport


DKW_ENDPOINT: Final[dict[str, str]] = {
    "stg": "https://dkw.dev.cloudnativedays.jp/query",
    "prd": "https://dkw.cloudnativedays.jp/query",
}


class Command:
    """CFP Vote Counter CLI"""

    def generate_csv(self, event: str, env: Literal["stg", "prd"] = "prd"):
        """
        Generate transformed CFP vote csv

        :param event: Event Abbreviation (e.g. cndt2022)
        :param env: Environment (stg|prd)
        """

        transport = AIOHTTPTransport(url=DKW_ENDPOINT[env])
        client = Client(transport=transport, fetch_schema_from_transport=True)

        query = gql(
        """
        query getVoteCounts($confName: ConfName!){
          voteCounts(confName: $confName) {
            talkId
            count
          }
        }
        """)

        data = client.execute(query, {"confName": event})
        df = (
            pd.DataFrame(data["voteCounts"])
            .sort_values(["count", "talkId"], ascending=[False, True])
            .reset_index(drop=True)
        )
        print(df.to_csv())


if __name__ == "__main__":
    fire.Fire(Command)
