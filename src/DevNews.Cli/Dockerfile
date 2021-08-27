FROM mcr.microsoft.com/dotnet/aspnet:5.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:5.0.400 AS build
WORKDIR /src
COPY . .
RUN dotnet tool restore
ENV PATH="/root/.dotnet/tools:${PATH}"

WORKDIR /src/src/DevNews.Cli
RUN dotnet paket restore
RUN dotnet restore -nowarn:msb3202,nu1503
RUN dotnet build -c Release -o /app

FROM build as unittests
WORKDIR /src/tests/DevNews.Core.UnitTests
RUN dotnet restore -nowarn:msb3202,nu1503
RUN dotnet run


FROM build as notificationsTests
WORKDIR /src/tests/DevNews.Infrastructure.Notifications.UnitTests
RUN dotnet restore -nowarn:msb3202,nu1503
RUN dotnet run

FROM build as parserTests
WORKDIR /src/tests/DevNews.Infrastructure.Notifications.UnitTests
RUN dotnet restore -nowarn:msb3202,nu1503
RUN dotnet run

FROM build AS publish
RUN dotnet publish --no-restore -c Release -o /app

FROM base AS final
WORKDIR /app
COPY --from=publish /app .
ENTRYPOINT ["dotnet", "DevNews.Cli.dll" ]
