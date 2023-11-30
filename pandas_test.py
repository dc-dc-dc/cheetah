import pandas as pd
import yfinance as yf

def download(symbol: str, start_date: str, end_date: str, time_frame: str):
    return yf.Ticker(symbol).history(start=start_date, end=end_date, interval=time_frame)

def sma(data: pd.DataFrame, window: int):
    data[f"sma{window}"] = data["Close"].rolling(window=window).mean()

def macd(data: pd.DataFrame):
    data["ema26"] = data["Close"].ewm(span=26).mean()
    data["ema12"] = data["Close"].ewm(span=12).mean()
    data["macd"] = data["ema12"] - data["ema26"]
    data["signal"] = data["macd"].ewm(span=9).mean()
    data["hist"] = data["macd"] - data["signal"]

if __name__ == "__main__":
    data = download("AAPL", "2022-11-29", "2023-11-29", "1d")
    sma(data, 20)
    macd(data)
    data.to_csv("testing_data.csv")