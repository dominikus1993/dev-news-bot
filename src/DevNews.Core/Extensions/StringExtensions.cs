using System;

namespace DevNews.Core.Extensions
{
    public static class StringExtensions
    {
        public static string TrimEntersAndSpaces(this string text)
        {
            return text.Trim().Replace(Environment.NewLine, string.Empty);
        } 
    }
}