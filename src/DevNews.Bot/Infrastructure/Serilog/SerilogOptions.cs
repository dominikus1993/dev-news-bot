namespace DevNews.DiscordBot.Infrastructure.Serilog
{
    public class SerilogOptions
    {
        public bool ConsoleEnabled { get; set; } = true;
        public string MinimumLevel { get; set; } = "Information";
        public string Format { get; set; } = "compact";
    }
}