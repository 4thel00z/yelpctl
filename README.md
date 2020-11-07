# yelpctl

Extract a subset of the [yelp academic dataset](https://www.yelp.com/dataset). Filter via bounding box.

## Example usage

This demonstrates how you can extract a subset of the dataset:

```bash
go run cmd/yelpctl/main.go -bbox "47.2701114,55.099161,5.8663153,15.0419319" -path assets/yelp_academic_dataset_business.json
```
## License

This project is licensed under the GPL-3 license.
