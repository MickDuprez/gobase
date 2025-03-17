*** GoBase
GoBase is a very slim micro framework that focuses on creating a base application to handle the day to day work such as basic Auth, database and Session management etc.

The goal is to provide a simple system to allow the developer to add 'features' to the base without interfering with other 'features', sort of like having your own lane or silo for each feature. Each feature, at the risk of some code duplication, handles its own database interface logic and html pages etc. This enables features to be developed in isolation, you can think of features like 'plugins' in this regard.

The UI core relies on a base.html layout page where main menu items are added as provided by the feature, the feature then may implement a more feature specific UI layout as its base which then renders the different feature page templates. The CSS is Bootstrap v5+ but this can be overwritten easily enough as the base.html is very simple. 

The main goal is to keep the 'core' as bare-bones as possible and if a feature requires more it can implement it as needed.

The project is in its very early stages and subject to frequent changes while I use it in some production projects and tweak or add functionality as I go.

** TODO 
[] Add example of using session data storage
