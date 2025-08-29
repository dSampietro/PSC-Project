import matplotlib.pyplot as plt
import pandas as pd

df = pd.read_csv("analysis.csv", header=0)
df["speedup"] = df["time_seq"] / df["time_par"]
df["unit_time_seq"] = df["time_seq"] / df["num_sent"]
df["unit_time_par"] = df["time_par"] / df["num_sent"]


print(df)

df.plot(x="depth", y=["time_seq", "time_par"], marker="o")

plt.title("Sequential vs parallel time")
plt.xlabel("depth")
plt.ylabel("time[s]")
plt.grid(True)



df.plot(x="depth", y="speedup", marker="o")

plt.title("Speedup")
plt.xlabel("depth")
plt.ylabel("speedup")
plt.grid(True)


df.plot(x="depth", y="num_sent", marker="o")

plt.title("Number of generated sentences")
plt.xlabel("depth")
plt.ylabel("#sentences")
plt.grid(True)



plt.show()
