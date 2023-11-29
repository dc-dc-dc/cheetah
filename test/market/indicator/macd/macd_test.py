import os
import yfinance as yf

def run():
    # Symbol
    symbol = os.getenv("SYMBOL", "AAPL")
    start_date = os.getenv("START_TIME", "2022-11-29")
    end_date = os.getenv("END_TIME", "2023-11-29")
    time_frame = os.getenv("TIMEFRAME", "1d")
    data = yf.Ticker(symbol).history(start=start_date, end=end_date, interval=time_frame)
    data["ema26"] = data["Close"].ewm(span=26).mean()
    data["ema12"] = data["Close"].ewm(span=12).mean()
    data["macd"] = data["ema12"] - data["ema26"]
    data["signal"] = data["macd"].ewm(span=9).mean()
    data["hist"] = data["macd"] - data["signal"]
    data.to_csv("macd_test_py.csv")

if __name__ == "__main__":
    run()