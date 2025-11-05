export const Testimonial = () => {
  return (
    <section className="mt-22">
      <h2 className="sr-only">Testimonials</h2>
      <div className="grid grid-cols-1 grid-rows-[auto_1fr] gap-x-10 gap-y-10 lg:grid-cols-2 lg:gap-y-5 border-y">
        <figure className="group row-span-2 grid mx-auto max-w-3xl gap-y-5 lg:grid-rows-subgrid  lg:first:border-r lg:last:border-l">
          <blockquote className="mx-auto flex items-center px-8 py-2 text-xl/9 font-medium tracking-tight sm:px-16 sm:text-2xl/10 lg:group-first:border-b">
            <p className="relative before:pointer-events-none before:absolute before:top-4 before:-left-6 before:text-[6rem] before:text-gray-950/10 before:content-['“'] sm:before:-left-8 lg:before:text-[8rem] dark:before:text-white/10">
              Rapida helped us add Voice AI capabilities without touching our
              core infrastructure. What used to take months of engineering now
              happens in weeks, and our customers experience it as a native
              feature.
            </p>
          </blockquote>
          <figcaption className="grid grid-cols-[max-content_1fr] gap-6 px-8 py-2 border-y sm:border-b-0 sm:px-16">
            <div className="text-sm/7">
              <p className="font-medium">Priya Menon</p>
              <p className="text-gray-600 dark:text-gray-400">
                Head of Product, Mid-Market CPaaS Provider
              </p>
            </div>
          </figcaption>
        </figure>
        <figure className="group row-span-2 grid mx-auto max-w-3xl gap-y-5 lg:grid-rows-subgrid  lg:first:border-r lg:last:border-l">
          <blockquote className="mx-auto flex items-center px-8 py-2 text-xl/9 font-medium tracking-tight border-b sm:px-16 sm:text-2xl/10 lg:group-first:border-b">
            <p className="relative before:pointer-events-none before:absolute before:top-4 before:-left-6 before:text-[6rem] before:text-gray-950/10 before:content-['“'] sm:before:-left-8 lg:before:text-[8rem] dark:before:text-white/10">
              Integrating AI workflows used to be painful and brittle. With
              Rapida, our latency stayed predictable, and our teams could focus
              on delivering new customer experiences instead of building
              pipelines.
            </p>
          </blockquote>
          <figcaption className="grid grid-cols-[max-content_1fr] gap-6 px-8 py-2 border-t sm:px-16">
            <div className="text-sm/7">
              <p className="font-medium">Daniel Koh</p>
              <p className="text-gray-600 dark:text-gray-400">
                CTO, Cloud Communications Platform
              </p>
            </div>
          </figcaption>
        </figure>
      </div>
    </section>
  );
};
