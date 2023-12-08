from setuptools import setup, find_packages


with open("requirements.txt") as f:
    requirements = f.read().splitlines()


setup(
    name="pears_dht",
    version="0.0.1",
    description="DHT network setup for PeARSearch",
    author="Nandaja Varma",
    author_email="nandaja.varma@gmail.com",
    url="https://github.com/PeARSearch/PeARS-dht",
    packages=find_packages(exclude=['test']),
    install_requires=requirements,
)
