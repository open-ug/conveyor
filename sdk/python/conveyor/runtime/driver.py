
class Driver:

    def __init__(self, name, resources: list[str]):
        self.name = name
        self.resources = resources

    def reconcile(self, payload: str, event: str, run_id: str, driver_logger):
        pass
