namespace cinema.log.server.Models.DTOs;

public abstract class Mapper<TSource, TDestination>
{
    /// <summary>
    /// Maps object of type Ts onto object of type Td.
    /// Looks at property names and types so very important that the names match and the types
    /// </summary>
    /// <param name="source"></param>
    /// <returns></returns>
    /// <exception cref="NullReferenceException"></exception>
    public static TDestination Map(TSource source)
    {
        if (source == null) throw new NullReferenceException(nameof(source));
        var T = Activator.CreateInstance<TDestination>();
        
        var sourceProperties = source.GetType().GetProperties();
        var targetProperties = T.GetType().GetProperties();
        
        foreach (var prop in sourceProperties)
        {
            var name = prop.Name;
            var targetProp = targetProperties.FirstOrDefault(t => t.Name == name);
            
            var value = prop.GetValue(source);
            if (targetProp != null && targetProp.PropertyType.IsAssignableFrom(prop.PropertyType))
            {
                targetProp.SetValue(T, value);
            }
        }
        return T;
    }
}