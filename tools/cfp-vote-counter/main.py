#!/usr/bin/env python3
from typing import Final, Literal, Optional

import fire
import pandas as pd
import matplotlib.pyplot as plt

from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport

VOTING_TERM: Final[dict[str, tuple[str,str]]] = {
    "cicd2023": ('2023-01-01T00:00:00+09:00', '2023-01-25T18:00:00+09:00'),
    "cndf2023": ('2023-05-02T00:00:00+09:00', '2023-06-19T23:59:00+09:00'),
    "cndt2023": ('2023-09-01T00:00:00+09:00', '2023-10-23T18:00:00+09:00'),
    "cnds2024": ('2024-03-08T00:00:00+09:00', '2024-04-09T19:00:00+09:00'),
}


DKW_ENDPOINT: Final[dict[str, str]] = {
    "dev": "http://localhost:8080/query",
    "stg": "https://dkw.dev.cloudnativedays.jp/query",
    "prd": "https://dkw.cloudnativedays.jp/query",
}

class Command:
    """CFP Vote Counter CLI"""

    def generate_csv(self, event: str, env: Literal["dev","stg", "prd"] = "prd"):
        """
        Generate transformed CFP vote csv

        :param event: Event Abbreviation (e.g. cndt2023)
        :param env: Environment (dev|stg|prd)
        """

        transport = AIOHTTPTransport(url=DKW_ENDPOINT[env])
        client = Client(transport=transport, fetch_schema_from_transport=True)

        query = gql(
        """
        query getVoteCounts($confName: ConfName!, $votingTerm: VotingTerm!){
          voteCounts(confName: $confName, votingTerm: $votingTerm) {
            talkId
            count
          }
        }
        """)

        if event not in VOTING_TERM:
            raise ValueError(f"event not exist: name={event}")

        data = client.execute(query, {
          "confName": event,
          "votingTerm":{
            "start": VOTING_TERM[event][0],
            "end":VOTING_TERM[event][1],
          }
        })
        print(data)
        df = (
            pd.DataFrame(data["voteCounts"])
            .sort_values(["count", "talkId"], ascending=[False, True])
            .reset_index(drop=True)
        )
        print(df.to_csv())


if __name__ == "__main__":
    fire.Fire(Command)
