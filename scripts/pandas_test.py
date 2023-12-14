import pandas as pd
import yfinance as yf

def download(symbol: str, start_date: str, end_date: str, time_frame: str):
    return yf.Ticker(symbol).history(start=start_date, end=end_date, interval=time_frame)

def rsi(data: pd.DataFrame, window: int):
    data["loss"] = data["Close"].diff(1).apply(lambda x: 0 if x > 0 else abs(x))
    data["gain"] = data["Close"].diff(1).apply(lambda x: 0 if x < 0 else x)
    data["avg_loss"] = data["loss"].rolling(window=window).mean()
    data["avg_gain"] = data["gain"].rolling(window=window).mean()
    data[f"rsi{window}"] = 100 - (100 / (1 + data["avg_gain"] / data["avg_loss"]))

def _min(data: pd.DataFrame, window: int):
    data[f"min{window}"] = data["Close"].rolling(window=window).min()

def _max(data: pd.DataFrame, window: int):
    data[f"max{window}"] = data["Close"].rolling(window=window).max()

def sma(data: pd.DataFrame, window: int):
    data[f"sma{window}"] = data["Close"].rolling(window=window).mean()

def ema(data: pd.DataFrame, window: int):
    data[f"ema{window}"] = data["Close"].ewm(span=window, adjust=False).mean()

def macd(data: pd.DataFrame):
    ema(data, 12)
    ema(data, 26)
    data["macd"] = data["ema12"] - data["ema26"]
    data["macd.signal"] = data["macd"].ewm(span=9, adjust=False).mean()
    data["macd.histogram"] = data["macd"] - data["macd.signal"]

if __name__ == "__main__":
    data = download("AAPL", "2022-11-29", "2023-11-29", "1d")
    _min(data, 20)
    _max(data, 20)
    rsi(data, 14)
    sma(data, 20)
    macd(data)
    data.to_csv("testing_data.csv")