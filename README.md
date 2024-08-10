<div align="center">
  <img src="public/photo.webp" alt="Repocate Image">
</div>

# Repocate: Your Code's Favorite Moving Company! 📦🚚

Ever wish you could pack up your entire dev environment and move it to a new machine faster than you can say "git clone"? Well, pack your bags (or don't, actually), because Repocate is here to do the heavy lifting for you!

Repocate is the turbocharged, no-hassle way to relocate your development setup. It's like having a team of coding movers who know exactly where to put your Node modules and won't ever lose your Go packages in transit.

### With Repocate, you can:

- Clone repos faster than a caffeinated developer types "npm install"
- Spin up dev environments quicker than you can say "Docker run"
- Jump between projects like a ninja hopping across rooftops

No more "It works on my machine" blues. No more spending half your day setting up a new environment. Just pure, unadulterated coding bliss.

So, whether you're a digital nomad hopping between coffee shops or a team lead onboarding new devs, Repocate is your one-way ticket to Productivity City. All aboard the express train to Efficient-ville!

Ready to make your code feel at home anywhere? Let's repocate! 🚀


## Features

- **Containerized Development Environment**: Repocate provides a consistent, reproducible, and isolated development environment inside Docker containers. It includes all necessary dependencies and configurations, ensuring that developers can get started quickly without worrying about local environment setup.

- **Dynamic Port and Volume Management**: Repocate intelligently manages dynamic port assignments and Docker volumes to avoid conflicts, making it easier to run multiple instances simultaneously.

- **Version Control Integration**: The tool integrates seamlessly with Git, allowing automatic versioning and tagging of development containers. This ensures that every significant change is captured and can be easily reverted if necessary.

- **Custom Plugin Support**: Extend the environment with custom Zsh plugins, ensuring that your development shell is tailored to your specific needs.

- **Snapshot and Rollback**: Create Git snapshots before major changes, and easily roll back to previous states in case of issues, ensuring a robust development process.

- **Flexible Configuration**: Repocate allows users to customize their development environment through `.zshrc` and Neovim configuration files, enabling a highly personalized development experience.
- 

## Usage

After setting up Repocate, you can start using it by cloning repositories into containers. Here's a quick guide:

1. **Create a Development Container:**

   ```sh
   repocate create <repo-url>
   ```

   This command will clone the specified repository into a containerized environment.

2. **Enter the Development Container:**

   ```sh
   repocate enter <repo-url> or <repo_name>
   ```

   Once the container is created, you can enter it using this command.

3. **Stop the Container:**

   ```sh
   repocate stop <repo-url>
   ```

   Stops the running development container.

4. **Rebuild the Container:**

   ```sh
   repocate rebuild <repo-url>
   ```

   Rebuilds the container, applying any updates or changes.

5. **List All Containers:**

   ```sh
   repocate list
   ```

   Lists all Repocate-managed containers with their current status.

6. **Advanced Usage:**

   - **Snapshot:** `repocate snapshot` to create a Git snapshot.
   - **Rollback:** `repocate rollback` to revert to the last known good commit.
   - **Volume Management:** Dynamically create and manage Docker volumes with `repocate volume`.

## Getting Started

1. **Install Repocate:**

   Clone the repository and run the following commands:

   ```sh
   git clone https://github.com/cdaprod/repocate.git
   cd repocate
   make install
   ```

   This will install the Repocate script and set up necessary configurations.

2. **Configure Your Environment:**

   Customize your `.zshrc` and Neovim settings to tailor the environment to your needs. Configuration files are located in `~/.config/repocate/`.

3. **Build Your Docker Image:**

   You can build and push Docker images using the provided GitHub Actions workflows, which handle dynamic port management, caching, and versioning.

## Contributing

We welcome contributions to improve Repocate! To contribute:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a Pull Request.

## License

Repocate is licensed under the MIT License. See the `LICENSE` file for more information.

## Support

If you encounter any issues or have questions, feel free to open an issue on GitHub or reach out to the maintainer at cdaprod@contact.com.

