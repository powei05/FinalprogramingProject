# app.R
# An R Shiny app to visualize ecosystem simulation results.

# --- 1. Load necessary libraries ---
# Make sure you have these installed: install.packages(c("shiny", "ggplot2", "plotly", "tidyr"))
library(shiny)
library(ggplot2)
library(plotly)
library(tidyr)

# --- 2. Define the User Interface (UI) ---
ui <- fluidPage(
    
    # Application title
    titlePanel("Ecosystem Population Dynamics"),
    
    # Sidebar with a checkbox to select species
    sidebarLayout(
        sidebarPanel(
            h4("Display Options"),
            # Checkboxes to toggle species visibility on the plot
            checkboxGroupInput("species_to_show", 
                               "Select Species to Display:",
                               choices = c("rabbit", "sheep", "deer", "wolf", "human"),
                               selected = c("rabbit", "wolf")), # Initially select rabbit and wolf
            
            hr(), # Horizontal line
            
            helpText("Data is read from 'population_log.csv' in the same directory.")
        ),
        
        # Main panel to display the plot
        mainPanel(
            # Using plotlyOutput for an interactive plot
            plotlyOutput("populationPlot", height = "600px")
        )
    )
)

# --- 3. Define the Server Logic ---
server <- function(input, output) {
    
    # Reactive expression to read and prepare the data
    # This will re-run automatically if the file changes (e.g., after a new simulation)
    population_data <- reactiveFileReader(1000, NULL, 'population_log.csv', read.csv)
    
    # Render the interactive plot
    output$populationPlot <- renderPlotly({
        
        # Ensure data is loaded
        req(population_data())
        
        # Reshape data from wide to long format for ggplot
        data_long <- pivot_longer(population_data(), 
                                  cols = -Generation, 
                                  names_to = "Species", 
                                  values_to = "Population")
        
        # Filter data based on user's checkbox selection
        filtered_data <- subset(data_long, Species %in% input$species_to_show)
        
        # Create the ggplot object
        p <- ggplot(filtered_data, aes(x = Generation, y = Population, color = Species)) +
            geom_line(size = 1) +
            labs(
                title = "Population over Generations",
                x = "Generation (Time)",
                y = "Population Count"
            ) +
            scale_color_manual(values = c(
                "rabbit" = "lightpink",
                "sheep"  = "lightgreen",
                "deer"   = "lightblue",
                "wolf"   = "darkgray",
                "human"  = "gold"
            )) +
            theme_minimal() +
            theme(legend.position = "bottom")
        
        # Convert the ggplot object to a plotly object for interactivity
        ggplotly(p)
    })
}

# --- 4. Run the application ---
shinyApp(ui = ui, server = server)

